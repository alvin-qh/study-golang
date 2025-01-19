'strict';

export class Clock {
  constructor() {
    this._date = new Date();
  }

  static _padNumber(num, pad = 2) {
    return `${num}`.padStart(pad, '0');
  }

  getYear() {
    return Clock._padNumber(this._date.getFullYear(), 4);
  }

  getMonth() {
    return Clock._padNumber(this._date.getMonth() + 1);
  }

  getDate() {
    return Clock._padNumber(this._date.getDate());
  }

  getHours() {
    return Clock._padNumber(this._date.getHours());
  }

  getMinutes() {
    return Clock._padNumber(this._date.getMinutes());
  }

  getSeconds() {
    return Clock._padNumber(this._date.getSeconds());
  }
}
