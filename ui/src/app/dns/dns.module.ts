import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {DnsComponent} from './dns.component';
import {CoreModule} from '../core/core.module';


@NgModule({
  declarations: [DnsComponent],
  imports: [
    CommonModule,
    CoreModule
  ]
})
export class DnsModule {
}
