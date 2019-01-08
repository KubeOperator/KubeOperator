import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ConfigComponent} from './config.component';
import {NodeConfigComponent} from './node-config/node-config.component';
import {ClusterConfigComponent} from './cluster-config/cluster-config.component';
import {CoreModule} from '../core/core.module';
import {ConfigService} from './config.service';
import {ConfigControlService} from './cluster-config/config-control.service';
import { ConfigItemComponent } from './cluster-config/config-item/config-item.component';

@NgModule({
  declarations: [ConfigComponent, NodeConfigComponent, ClusterConfigComponent, ConfigItemComponent],
  imports: [
    CommonModule,
    CoreModule
  ],
  providers: [ConfigService, ConfigControlService]
})
export class ConfigModule {
}
