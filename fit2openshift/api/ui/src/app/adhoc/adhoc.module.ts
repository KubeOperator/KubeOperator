import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CoreModule } from '../core/core.module';
import { SharedModule } from '../shared/shared.module';
import { AdhocComponent } from './adhoc/adhoc.component';

@NgModule({
  imports: [
    CommonModule,
    CoreModule,
    SharedModule,
  ],
  declarations: [AdhocComponent]
})
export class AdhocModule { }
