import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {DeployComponent} from './deploy.component';
import {TermComponent} from './term/term.component';
import {OperaterComponent} from './operater/operater.component';
import {CoreModule} from '../core/core.module';
import {ProgressComponent} from './progress/progress.component';

@NgModule({
  declarations: [DeployComponent, TermComponent, OperaterComponent, ProgressComponent],
  imports: [
    CommonModule,
    CoreModule
  ]
})
export class DeployModule {
}
