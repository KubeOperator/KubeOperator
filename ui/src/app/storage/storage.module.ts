import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {StorageComponent} from './storage.component';
import {CoreModule} from '../core/core.module';


@NgModule({
  declarations: [StorageComponent],
  imports: [
    CommonModule,
    CoreModule
  ]
})
export class StorageModule {
}
