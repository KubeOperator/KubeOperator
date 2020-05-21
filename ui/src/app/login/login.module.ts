import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LoginComponent } from './login.component';
import {CoreModule} from '../core/core.module';



@NgModule({
  declarations: [LoginComponent],
    imports: [
        CommonModule,
        CoreModule
    ]
})
export class LoginModule { }
