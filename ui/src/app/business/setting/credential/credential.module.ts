import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CredentialDeleteComponent } from './credential-delete/credential-delete.component';
import {CoreModule} from '../../../core/core.module';



@NgModule({
  declarations: [CredentialDeleteComponent],
  exports: [
    CredentialDeleteComponent
  ],
  imports: [
    CommonModule,
    CoreModule
  ]
})
export class CredentialModule { }
