import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LoginComponent } from './login.component';
import {CoreModule} from '../core/core.module';
import {SharedModule} from '../shared/shared.module';
import { ForgotPasswordComponent } from './forgot-password/forgot-password.component';



@NgModule({
  declarations: [LoginComponent, ForgotPasswordComponent],
    imports: [
        CommonModule,
        CoreModule,
        SharedModule
    ]
})
export class LoginModule { }
