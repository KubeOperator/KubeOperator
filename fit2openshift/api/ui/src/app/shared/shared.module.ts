import { NgModule } from '@angular/core';
import { CoreModule } from '../core/core.module';
import { PageNotFoundComponent } from './not-found/not-found.component';
import { AuthService } from './auth.service';

@NgModule({
  imports: [
    CoreModule,
  ],
  exports: [
    // CoreModule
  ],
  declarations: [PageNotFoundComponent],
  providers: [AuthService]

})
export class SharedModule { }
