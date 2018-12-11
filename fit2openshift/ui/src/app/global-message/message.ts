import {AlertType} from '../shared/shared.const';

export class Message {
  statusCode: number;
  message: string;
  alertType: AlertType;
  isAppLevel = false;

  get type(): string {
    switch (this.alertType) {
      case AlertType.DANGER:
        return 'alert-danger';
      case AlertType.INFO:
        return 'alert-info   ';
      case AlertType.SUCCESS:
        return 'alert-success';
      case AlertType.WARNING:
        return 'alert-warning';
      default:
        return 'alert-warning';
    }
  }

  static newMessage(statusCode: number, message: string, alertType: AlertType): Message {
    const m = new Message();
    m.statusCode = statusCode;
    m.message = message;
    m.alertType = alertType;
    return m;
  }

  toString(): string {
    return 'Message with statusCode:' + this.statusCode +
      ', message:' + this.message +
      ', alert type:' + this.type;
  }
}
