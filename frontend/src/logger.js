import { createLogger, format, transports } from 'winston'
import { combine, timestamp, printf } from 'logform/format'

const stackframe = format((info, opts) => {
  const stack = new Error().stack
  
  const stackRows = stack.split('\n').pop()
  const [sourceFile, lineNumber] = stackRows.split(/\s+/).pop().split(':')
  const functionName = stackRows.match(/at\s(\S*)\s/)

  if (sourceFile) {
    info['sourceFile'] = sourceFile.replace("(", "")
  }

  if (lineNumber) {
    info['lineNumber'] = lineNumber
  }
  
  if (functionName) {
    info['functionName'] = functionName[1].trim()
  }

  return info
})

const stackdriverFormat = printf(({sourceFile, lineNumber, functionName, level, message, timestamp}) => { 
  
  const info = Object.assign({}, {
    'severity': level,
    'message': message,
    'timestamp': timestamp,
    'file': sourceFile,
    'line': lineNumber,
  })

  if (functionName) {
    info['function'] = functionName
  }

  return JSON.stringify(info)
})

const stackdriverLogger = createLogger({
  format: combine(
    timestamp(),
    stackframe(),
    stackdriverFormat
  ),
  transports: [new transports.Console()]
})

export default stackdriverLogger
