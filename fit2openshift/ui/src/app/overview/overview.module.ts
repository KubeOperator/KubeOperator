import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { OverviewComponent } from './overview.component';
import { ProgressComponent } from './progress/progress.component';
import { CharsComponent } from './chars/chars.component';
import {CoreModule} from '../core/core.module';
import { DescribeComponent } from './describe/describe.component';
import { TermComponent } from './term/term.component';
import { OperaterComponent } from './operater/operater.component';

@NgModule({
  declarations: [OverviewComponent, ProgressComponent, CharsComponent, DescribeComponent, TermComponent, OperaterComponent],
  imports: [
    CommonModule,
    CoreModule
  ]
})
export class OverviewModule { }
