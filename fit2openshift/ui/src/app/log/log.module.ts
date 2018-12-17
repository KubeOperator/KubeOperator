import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {LogComponent} from './log.component';
import {LogListComponent} from './log-list/log-list.component';
import {CoreModule} from '../core/core.module';
import {LogService} from './log.service';

@NgModule({
  declarations: [LogComponent, LogListComponent],
  imports: [
    CommonModule,
    CoreModule
  ],
  providers: [LogService]
})
export class LogModule {
}
