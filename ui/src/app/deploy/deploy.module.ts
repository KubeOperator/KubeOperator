import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {DeployComponent} from './deploy.component';
import {TermComponent} from './component/term/term.component';
import {CoreModule} from '../core/core.module';
import {ProgressComponent} from './component/progress/progress.component';

@NgModule({
  declarations: [DeployComponent, TermComponent, ProgressComponent],
  imports: [
    CommonModule,
    CoreModule
  ]
})
export class DeployModule {
}
