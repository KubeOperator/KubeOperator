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
import {LogModule} from './log/log.module';
import {HostModule} from './host/host.module';
import {DeployModule} from './deploy/deploy.module';
import {SettingModule} from './setting/setting.module';
import {CredentialModule} from './credential/credential.module';
import {RegionModule} from './region/region.module';
import {ZoneModule} from './zone/zone.module';
import {PlanModule} from './plan/plan.module';
import {F5BigIpModule} from './f5-big-ip/f5-big-ip.module';
import {ClusterHealthComponent} from './cluster-health/cluster-health.component';
import {DeployPlanComponent} from './deploy-plan/deploy-plan.component';
import {ApplicationComponent} from './application/application.component';
import {ClusterBackupComponent} from './cluster-backup/cluster-backup.component';
import {ClusterBackupModule} from './cluster-backup/cluster-backup.module';
import {SharedModule} from './shared/shared.module';
import {NfsModule} from './nfs/nfs.module';

@NgModule({
  declarations: [
    AppComponent,
    DeployPlanComponent,
    ApplicationComponent,
    ClusterHealthComponent,
    ClusterBackupComponent,
  ],
  imports: [
    CredentialModule,
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
    RegionModule,
    NodeModule,
    LogModule,
    HostModule,
    SettingModule,
    ZoneModule,
    PlanModule,
    F5BigIpModule,
    ClusterBackupModule,
    SharedModule,
  ],
  providers: [{provide: HTTP_INTERCEPTORS, useClass: InterceptorService, multi: true}],
  bootstrap: [AppComponent]
})
export class AppModule {
}
