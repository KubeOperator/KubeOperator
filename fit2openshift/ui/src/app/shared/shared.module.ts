import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {GlobalMessageComponent} from '../global-message/global-message.component';
import {GlobalMessageService} from '../global-message/global-message.service';
import {NotFoundComponent} from './not-found/not-found.component';
import {CoreModule} from '../core/core.module';
import {AuthUserActiveService} from './route/auth-user-active.service';
import {MessageHandlerService} from './message-handler/message-handler.service';

@NgModule({
  declarations: [GlobalMessageComponent, NotFoundComponent],
  imports: [
    CommonModule,
    CoreModule
  ], exports: [
    GlobalMessageComponent,
    CoreModule
  ], providers: [
    GlobalMessageService,
    AuthUserActiveService,
    MessageHandlerService,

  ]
})
export class SharedModule {
}
