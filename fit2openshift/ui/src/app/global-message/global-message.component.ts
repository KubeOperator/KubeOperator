import {Component, ElementRef, Input, OnDestroy, OnInit} from '@angular/core';
import {Message} from './message';
import {Subscription} from 'rxjs';
import {GlobalMessageService} from './global-message.service';
import {Router} from '@angular/router';
import {dismissInterval} from '../shared/shared.const';

@Component({
  selector: 'app-global-message',
  templateUrl: './global-message.component.html',
  styleUrls: ['./global-message.component.css']
})
export class GlobalMessageComponent implements OnInit, OnDestroy {

  @Input() isAppLevel: boolean;
  globalMessage: Message = new Message();
  globalMessageOpened: boolean;
  private messageText = '';
  timer: any = null;

  appLevelMsgSub: Subscription;
  msgSub: Subscription;
  clearSub: Subscription;

  constructor(private elementRef: ElementRef, private messageService: GlobalMessageService, private router: Router) {
  }

  ngOnInit() {
    if (this.isAppLevel) {
      this.appLevelMsgSub = this.messageService.appLevelAnnounced$.subscribe(message => {
        this.globalMessageOpened = true;
        this.globalMessage = message;
        this.messageText = message.message;
      });
    } else {
      this.msgSub = this.messageService.messageAnnounced$.subscribe(message => {
          this.globalMessageOpened = true;
          this.globalMessage = message;
          this.messageText = message.message;

          this.timer = setTimeout(() => this.onClose(), dismissInterval);

          setTimeout(() => {
            const nativeDom: any = this.elementRef.nativeElement;
            const queryDoms: any[] = nativeDom.getElementsByClassName('alert');
            if (queryDoms && queryDoms.length > 0) {
              const hackDom: any = queryDoms[0];
              hackDom.className += ' alert-global alert-global-align';
            }
          }, 0);
        }
      );
    }
    this.clearSub = this.messageService.clearChan$.subscribe(() => {
      this.onClose();
    });
  }


  ngOnDestroy(): void {
    if (this.appLevelMsgSub) {
      this.appLevelMsgSub.unsubscribe();
    }

    if (this.msgSub) {
      this.msgSub.unsubscribe();
    }

    if (this.clearSub) {
      this.clearSub.unsubscribe();
    }
  }


  get message(): string {
    return this.messageText;
  }

  onClose() {
    if (this.timer) {
      clearTimeout(this.timer);
    }
    this.globalMessageOpened = false;
  }


}
