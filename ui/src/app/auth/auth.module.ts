import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {CoreModule} from '../core/core.module';
import {AuthComponent} from './auth.component';

@NgModule({
  declarations: [AuthComponent],
  imports: [
    CommonModule,
    CoreModule
  ]
})
export class AuthModule {
}
