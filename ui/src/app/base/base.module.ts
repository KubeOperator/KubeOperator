import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FooterComponent} from './footer/footer.component';
import {HeaderComponent} from './header/header.component';
import {NavigatorComponent} from './navigator/navigator.component';
import {ShellComponent} from './shell/shell.component';
import {CoreModule} from '../core/core.module';
import {PasswordComponent} from './header/components/password/password.component';

import {CommonAlertComponent} from './header/components/common-alert/common-alert.component';

@NgModule({
  declarations: [FooterComponent, HeaderComponent, NavigatorComponent, ShellComponent,
     PasswordComponent, CommonAlertComponent],
  imports: [
    CommonModule,
    CoreModule,
  ],
  exports: [NavigatorComponent]
})
export class BaseModule {
}
