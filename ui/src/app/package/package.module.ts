import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {PackageComponent} from './package.component';
import {PackageService} from './package.service';
import {SharedModule} from '../shared/shared.module';
import {PackageListComponent} from './package-list/package-list.component';
import {CoreModule} from '../core/core.module';

import { PackageDetailComponent } from './package-detail/package-detail.component';

@NgModule({
  declarations: [PackageComponent, PackageListComponent, PackageDetailComponent],
  imports: [
    CommonModule,
    SharedModule,
    CoreModule,
  ], providers: [PackageService]
})
export class PackageModule {
}
