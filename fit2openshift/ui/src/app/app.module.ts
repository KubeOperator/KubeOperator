import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {BaseModule} from './base/base.module';
import {AccountModule} from './account/account.module';
import {InterceptorService} from './shared/interceptor.service';
import {HTTP_INTERCEPTORS} from '@angular/common/http';
import {PackageModule} from './package/package.module';
import {UserModule} from './user/user.module';
import {ClusterModule} from './cluster/cluster.module';
import {OverviewModule} from './overview/overview.module';
import {NodeModule} from './node/node.module';
import {ConfigModule} from './config/config.module';
import {MonitorModule} from './monitor/monitor.module';
import {LogModule} from './log/log.module';
import {TipModule} from './tip/tip.module';
import {HostModule} from './host/host.module';
import {HostFilterPipe} from './host/host-filter.pipe';
import {DeployModule} from './deploy/deploy.module';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    BaseModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    AccountModule,
    PackageModule,
    UserModule,
    ClusterModule,
    DeployModule,
    OverviewModule,
    NodeModule,
    ConfigModule,
    LogModule,
    MonitorModule,
    TipModule,
    HostModule,
  ],
  providers: [{provide: HTTP_INTERCEPTORS, useClass: InterceptorService, multi: true}],
  bootstrap: [AppComponent]
})
export class AppModule {
}
