import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {PackageComponent} from './package.component';
import {PackageService} from './package.service';
import {SharedModule} from '../shared/shared.module';
import {PackageListComponent} from './package-list/package-list.component';
import {CoreModule} from '../core/core.module';
import {TipModule} from '../tip/tip.module';

@NgModule({
  declarations: [PackageComponent, PackageListComponent],
  imports: [
    CommonModule,
    SharedModule,
    CoreModule,
    TipModule
  ], providers: [PackageService]
})
export class PackageModule {
}
