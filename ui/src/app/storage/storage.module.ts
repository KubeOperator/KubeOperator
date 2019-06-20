import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {StorageComponent} from './storage.component';
import {StorageListComponent} from './components/storage-list/storage-list.component';
import {StorageCreateComponent} from './components/storage-create/storage-create.component';
import {TipModule} from '../tip/tip.module';
import {CoreModule} from '../core/core.module';
import {SharedModule} from '../shared/shared.module';
import { StorageDetailComponent } from './components/storage-detail/storage-detail.component';

@NgModule({
  declarations: [StorageComponent, StorageListComponent, StorageCreateComponent, StorageDetailComponent],
  imports: [
    CommonModule,
    TipModule,
    CoreModule,
    SharedModule
  ]
})
export class StorageModule {
}
