import {Injectable} from '@angular/core';
import {GlobalMessageService} from '../../global-message/global-message.service';
import {SessionService} from '../session.service';
import {errorHandler} from '../../shared/shared.utils';
import {AlertType, httpStatusCode} from '../shared.const';


@Injectable()
export class MessageHandlerService {

  constructor(private msgService: GlobalMessageService, private sessionService: SessionService) {
  }

  public handleError(error: any | string): void {
    if (!error) {
      return;
    }
    const msg = errorHandler(error);

    if (!(error.statusCode || error.status)) {
      this.msgService.announceMessage(500, msg, AlertType.DANGER);
    } else {
      const code = error.statusCode || error.status;
      if (code === httpStatusCode.Unauthorized) {
        this.msgService.announceAppLevelMessage(code, msg, AlertType.DANGER);
        // Session is invalid now, clare session cache
        this.sessionService.clear();
      } else {
        this.msgService.announceMessage(code, msg, AlertType.DANGER);
      }
    }
  }

  public handleReadOnly(): void {
    this.msgService.announceAppLevelMessage(503, 'REPO_READ_ONLY', AlertType.WARNING);
  }

  public showError(message: string): void {
    this.msgService.announceMessage(500, message, AlertType.DANGER);
  }

  public showSuccess(message: string): void {
    if (message && message.trim() !== '') {
      this.msgService.announceMessage(200, message, AlertType.SUCCESS);
    }
  }

  public showInfo(message: string): void {
    if (message && message.trim() !== '') {
      this.msgService.announceMessage(200, message, AlertType.INFO);
    }
  }

  public showWarning(message: string): void {
    if (message && message.trim() !== '') {
      this.msgService.announceMessage(400, message, AlertType.WARNING);
    }
  }

  public clear(): void {
    this.msgService.clear();
  }

  public isAppLevel(error: any): boolean {
    return error && error.statusCode === httpStatusCode.Unauthorized;
  }

  public error(error: any): void {
    this.handleError(error);
  }

  public warning(warning: any): void {
    this.showWarning(warning);
  }

  public info(info: any): void {
    this.showSuccess(info);
  }

  public log(log: any): void {
    this.showInfo(log);
  }

}
