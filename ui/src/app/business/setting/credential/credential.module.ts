import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CredentialDeleteComponent } from './credential-delete/credential-delete.component';
import {CoreModule} from '../../../core/core.module';
import { CredentialEditComponent } from './credential-edit/credential-edit.component';



@NgModule({
  declarations: [CredentialDeleteComponent, CredentialEditComponent],
  exports: [
    CredentialDeleteComponent,
    CredentialEditComponent
  ],
  imports: [
    CommonModule,
    CoreModule
  ]
})
export class CredentialModule { }
