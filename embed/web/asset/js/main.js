'strict';

import { Clock } from './clock.js';

const $timer = document.querySelector('div.main .timer');

function getNow() {
  const now = new Clock();

  const date = `${now.getYear()}-${now.getMonth()}-${now.getDate()}`;
  const time = `${now.getHours()}:${now.getMinutes()}:${now.getSeconds()}`;
  return `${date} ${time}`;
}

$timer.innerHTML = getNow();

const timer = setInterval(() => {
  $timer.innerHTML = getNow();
}, 500);
