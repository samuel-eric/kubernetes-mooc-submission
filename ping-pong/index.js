const express = require('express');

const app = express();
const PORT = 3001;

let requestCounter = 0;

app.get('/pingpong', (req, res) => {
  requestCounter++;
  res.send(`pong ${requestCounter}`);
});

app.listen(PORT, () => {
  console.log(`Ping-pong server is running on port ${PORT}`);
});