const fs = require('fs');
const path = require('path');
const { v4: uuidv4 } = require('uuid');

const LOG_DIR = path.join(__dirname, 'log');
const LOG_FILE_PATH = path.join(LOG_DIR, 'app_log.txt');
const WRITE_INTERVAL_MS = 5000;
const randomString = uuidv4();

function logTimestamp() {
  const timestamp = new Date().toISOString();
  const logEntry = `${timestamp}: ${randomString}\n`;

  fs.appendFile(LOG_FILE_PATH, logEntry, (err) => {
    if (err) {
      console.error('Failed to write to log file:', err);
      return;
    }
  });
}

logTimestamp();
setInterval(logTimestamp, WRITE_INTERVAL_MS);
