import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LogComponent } from './log.component';
import { LogListComponent } from './log-list/log-list.component';
import {CoreModule} from '../core/core.module';

@NgModule({
  declarations: [LogComponent, LogListComponent],
  imports: [
    CommonModule,
    CoreModule
  ]
})
export class LogModule { }
