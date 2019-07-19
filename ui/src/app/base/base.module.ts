import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FooterComponent} from './footer/footer.component';
import {HeaderComponent} from './header/header.component';
import {NavigatorComponent} from './navigator/navigator.component';
import {ShellComponent} from './shell/shell.component';
import {CoreModule} from '../core/core.module';
import {MessageComponent} from './message/message.component';
import {MessageService} from './message.service';
import {PasswordComponent} from './header/components/password/password.component';
import {TipModule} from '../tip/tip.module';

@NgModule({
  declarations: [FooterComponent, HeaderComponent, NavigatorComponent, ShellComponent, MessageComponent, PasswordComponent],
  imports: [
    CommonModule,
    CoreModule,
    TipModule
  ],
  exports: [NavigatorComponent]
})
export class BaseModule {
}
