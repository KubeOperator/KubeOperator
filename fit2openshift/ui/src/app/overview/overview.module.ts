import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {OverviewComponent} from './overview.component';
import {ProgressComponent} from '../deploy/progress/progress.component';
import {CharsComponent} from './chars/chars.component';
import {CoreModule} from '../core/core.module';
import {DescribeComponent} from './describe/describe.component';
import {TermComponent} from '../deploy/term/term.component';
import {OperaterComponent} from '../deploy/operater/operater.component';

@NgModule({
  declarations: [OverviewComponent, CharsComponent, DescribeComponent],
  imports: [
    CommonModule,
    CoreModule
  ]
})
export class OverviewModule {
}
