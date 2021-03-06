package log

var modules []*ModuleLogger

// ModuleLogger is the logger used for getting the specific modules.
type ModuleLogger struct {
	Name           string
	logger         LeveledLogger
	isDebugLeveled bool
	isLevelSetter  bool

	levelSetter  LevelSetter
	debugLeveled DebugLeveledLogger

	currentLevel Level
}

// NewModuleLogger creates new module logger for given 'name' of the module and an optional 'logger'.
func NewModuleLogger(name string, moduleLogger ...LeveledLogger) *ModuleLogger {
	mLogger := &ModuleLogger{
		Name: name,
	}
	modules = append(modules, mLogger)

	if len(moduleLogger) > 0 {
		mLogger.logger = moduleLogger[0]
		mLogger.debugLeveled, mLogger.isDebugLeveled = moduleLogger[0].(DebugLeveledLogger)

		Debug2f("Module Logger: '%s' set from the argument", name)

		mLogger.initializeLogger()
		depthGetter, isDepthGetter := mLogger.logger.(OutputDepthGetter)
		if isDepthGetter {
			depthSetter, isDepthSetter := mLogger.logger.(OutputDepthSetter)
			if isDepthSetter {
				depthSetter.SetOutputDepth(depthGetter.GetOutputDepth() + 1)
			}
		}
	} else if sub, ok := defaultLogger.(SubLogger); ok {
		Infof("Module Logger: '%s' created as a SubLogger", name)
		mLogger.logger = sub.SubLogger()
		mLogger.initializeLogger()
	} else {
		Infof("Module Logger: '%s' as a wrapper over default logger", name)
		mLogger.currentLevel = currentLevel
		return mLogger
	}
	return mLogger
}

func (m *ModuleLogger) initializeLogger() {
	if m.logger == nil {
		return
	}

	m.debugLeveled, m.isDebugLeveled = m.logger.(DebugLeveledLogger)
	lGetter, ok := m.logger.(LevelGetter)
	if ok {
		m.currentLevel = lGetter.GetLevel()
	} else {
		m.currentLevel = currentLevel
	}
	m.levelSetter, m.isLevelSetter = m.logger.(LevelSetter)
}

// Level gets the module logger level.
func (m *ModuleLogger) Level() Level {
	if m.logger != nil {
		return m.currentLevel
	}
	return currentLevel
}

// SetLevel sets the moduleLogger level.
func (m *ModuleLogger) SetLevel(level Level) {
	Debugf("Setting Module: '%s' Logger level to: '%s'", m.Name, level)
	m.currentLevel = level
	if m.isLevelSetter {
		m.levelSetter.SetLevel(level)
	}
}

// Debug3f writes the formatted debug3 log.
func (m *ModuleLogger) Debug3f(format string, args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelDebug3 {
			return
		}
	}
	format = m.name() + " " + format
	if m.logger != nil {
		if !m.isDebugLeveled {
			m.logger.Debugf(format, args...)
		} else {
			m.debugLeveled.Debug3f(format, args...)
		}
	} else {
		Debug3f(format, args...)
	}
}

// Debug2f writes the formatted debug2 log.
func (m *ModuleLogger) Debug2f(format string, args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelDebug2 {
			return
		}
	}
	format = m.name() + " " + format
	if m.logger != nil {
		if !m.isDebugLeveled {
			m.logger.Debugf(format, args...)
		} else {
			m.debugLeveled.Debug2f(format, args...)
		}
	} else {
		Debug2f(format, args...)
	}
}

// Debugf writes the formatted debug log.
func (m *ModuleLogger) Debugf(format string, args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelDebug {
			return
		}
	}
	format = m.name() + " " + format
	if m.logger != nil {
		m.logger.Debugf(format, args...)
	} else {
		Debugf(format, args...)
	}
}

// Infof writes the formatted info log.
func (m *ModuleLogger) Infof(format string, args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelInfo {
			return
		}
	}
	format = m.name() + " " + format
	if m.logger != nil {
		m.logger.Infof(format, args...)
	} else {
		Infof(format, args...)
	}
}

// Warningf writes the formatted warning log.
func (m *ModuleLogger) Warningf(format string, args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelWarning {
			return
		}
	}
	format = m.name() + " " + format
	if m.logger != nil {
		m.logger.Warningf(format, args...)
	} else {
		Warningf(format, args...)
	}
}

// Errorf writes the formatted error log.
func (m *ModuleLogger) Errorf(format string, args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelError {
			return
		}
	}
	format = m.name() + " " + format
	if m.logger != nil {
		m.logger.Errorf(format, args...)
	} else {
		Errorf(format, args...)
	}
}

// Fatalf writes the formatted fatal log.
func (m *ModuleLogger) Fatalf(format string, args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelCritical {
			return
		}
	}
	format = m.name() + " " + format
	if m.logger != nil {
		m.logger.Fatalf(format, args...)
	} else {
		Fatalf(format, args...)
	}
}

// Panicf writes the formatted panic log.
func (m *ModuleLogger) Panicf(format string, args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelCritical {
			return
		}
	}
	format = m.name() + " " + format
	if m.logger != nil {
		m.logger.Panicf(format, args...)
	} else {
		Panicf(format, args...)
	}
}

// Debug3 writes the debug3 level log.
func (m *ModuleLogger) Debug3(args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelDebug3 {
			return
		}
	}
	args = append([]interface{}{m.name(), " "}, args...)
	if m.logger != nil {
		if !m.isDebugLeveled {
			m.logger.Debug(args...)
		} else {
			m.debugLeveled.Debug3(args...)
		}
	} else {
		Debug3(args...)
	}
}

// Debug2 writes the debug2 level log.
func (m *ModuleLogger) Debug2(args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelDebug2 {
			return
		}
	}
	args = append([]interface{}{m.name(), " "}, args...)
	if m.logger != nil {
		if !m.isDebugLeveled {
			m.logger.Debug(args...)
		} else {
			m.debugLeveled.Debug2(args...)
		}
	} else {
		Debug2(args...)
	}
}

// Debug writes the debug level log.
func (m *ModuleLogger) Debug(args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelDebug {
			return
		}
	}
	args = append([]interface{}{m.name(), " "}, args...)
	if m.logger != nil {
		m.logger.Debug(args...)
	} else {
		Debug(args...)
	}
}

// Info writes the info level log.
func (m *ModuleLogger) Info(args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelInfo {
			return
		}
	}
	args = append([]interface{}{m.name(), " "}, args...)
	if m.logger != nil {
		m.logger.Info(args...)
	} else {
		Info(args...)
	}
}

// Warning writes the warning level log.
func (m *ModuleLogger) Warning(args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelWarning {
			return
		}
	}
	args = append([]interface{}{m.name(), " "}, args...)
	if m.logger != nil {
		m.logger.Warning(args...)
	} else {
		Warning(args...)
	}
}

// Error writes the error level log.
func (m *ModuleLogger) Error(args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelError {
			return
		}
	}
	args = append([]interface{}{m.name(), " "}, args...)
	if m.logger != nil {
		m.logger.Error(args...)
	} else {
		Error(args...)
	}
}

// Fatal writes the fatal level log.
func (m *ModuleLogger) Fatal(args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelCritical {
			return
		}
	}
	args = append([]interface{}{m.name(), " "}, args...)
	if m.logger != nil {
		m.logger.Fatal(args...)
	} else {
		Fatal(args...)
	}
}

// Panic writes the panic level log.
func (m *ModuleLogger) Panic(args ...interface{}) {
	if !m.isLevelSetter {
		if m.currentLevel != LevelUnknown && m.currentLevel > LevelCritical {
			return
		}
	}
	args = append([]interface{}{m.name(), " "}, args...)
	if m.logger != nil {
		m.logger.Panic(args...)
	} else {
		Panic(args...)
	}
}

func (m *ModuleLogger) name() string {
	return "[" + m.Name + "]"
}
