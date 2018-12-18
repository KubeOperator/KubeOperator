import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {NotFoundComponent} from './not-found/not-found.component';
import {CoreModule} from '../core/core.module';
import {AuthUserActiveService} from './route/auth-user-active.service';

@NgModule({
  declarations: [ NotFoundComponent],
  imports: [
    CommonModule,
    CoreModule
  ], exports: [
    CoreModule
  ], providers: [
    AuthUserActiveService,
  ]
})
export class SharedModule {
}
