import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { CoreModule } from '../core/core.module';
import { SharedModule } from '../shared/shared.module';

import { CeleryLogComponent } from './celery-log/celery-log.component';

@NgModule({
  imports: [
    CommonModule,
    CoreModule,
    SharedModule,
  ],
  declarations: [CeleryLogComponent],
  exports: [CeleryLogComponent]
})
export class CeleryModule { }
