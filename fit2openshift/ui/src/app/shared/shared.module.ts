import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {NotFoundComponent} from './not-found/not-found.component';
import {CoreModule} from '../core/core.module';
import {AuthUserActiveService} from './route/auth-user-active.service';
import {NullFilterPipe} from './pipe/null-filter.pipe';

@NgModule({
  declarations: [NotFoundComponent, NullFilterPipe],
  imports: [
    CommonModule,
    CoreModule
  ], exports: [
    CoreModule,
    NullFilterPipe
  ], providers: [
    AuthUserActiveService,
  ]
})
export class SharedModule {
}
