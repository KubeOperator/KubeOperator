import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CredentialDeleteComponent } from './credential-delete/credential-delete.component';



@NgModule({
  declarations: [CredentialDeleteComponent],
  exports: [
    CredentialDeleteComponent
  ],
  imports: [
    CommonModule
  ]
})
export class CredentialModule { }
