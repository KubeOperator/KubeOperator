import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { OverviewComponent } from './overview.component';
import { ProgressComponent } from './progress/progress.component';
import { CharsComponent } from './chars/chars.component';
import {CoreModule} from '../core/core.module';
import { DescribeComponent } from './describe/describe.component';

@NgModule({
  declarations: [OverviewComponent, ProgressComponent, CharsComponent, DescribeComponent],
  imports: [
    CommonModule,
    CoreModule
  ]
})
export class OverviewModule { }
