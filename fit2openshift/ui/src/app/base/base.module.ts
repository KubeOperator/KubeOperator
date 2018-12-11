import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FooterComponent} from './footer/footer.component';
import {HeaderComponent} from './header/header.component';
import {NavigatorComponent} from './navigator/navigator.component';
import {ShellComponent} from './shell/shell.component';
import {CoreModule} from '../core/core.module';

@NgModule({
  declarations: [FooterComponent, HeaderComponent, NavigatorComponent, ShellComponent],
  imports: [
    CommonModule,
    CoreModule
  ],
  exports: [NavigatorComponent]
})
export class BaseModule {
}
