const express = require('express');
const fs = require('fs');
const path = require('path');

const PORT = 3000;
const LOG_DIR = path.join(__dirname, 'log');
const LOG_FILE_PATH = path.join(LOG_DIR, 'app_log.txt');

const app = express();

app.get('/', (req, res) => {
  fs.readFile(LOG_FILE_PATH, 'utf8', (err, data) => {
    if (err) {
      if (err.code === 'ENOENT') {
        console.log('Log file not found yet. The writer may not have run.');
        return res.status(404).type('text/plain').send('Log file not found. Please wait for the writer service to create it.');
      }
      console.error('Failed to read log file:', err);
      return res.status(500).type('text/plain').send('Error reading log file.');
    }

    res.type('text/plain').send(data);
  });
});

app.listen(PORT, () => {
  console.log(`Listening for requests on http://localhost:${PORT}`);
  console.log(`Serving content from: ${LOG_FILE_PATH}`);
});