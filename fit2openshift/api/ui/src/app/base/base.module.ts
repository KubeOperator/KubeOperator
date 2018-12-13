import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SharedModule } from '../shared/shared.module';
import { CoreModule } from '../core/core.module';
import { ShellComponent } from './shell/shell.component';
import { NavigatorComponent } from './navigator/navigator.component';
import { HeaderComponent } from './header/header.component';

@NgModule({
  imports: [
    CommonModule,
    CoreModule,
    SharedModule,
  ],
  exports: [
    NavigatorComponent
  ],
  declarations: [ShellComponent, NavigatorComponent, HeaderComponent]
})
export class BaseModule { }
