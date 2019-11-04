import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {ShellComponent} from './base/shell/shell.component';
import {NotFoundComponent} from './shared/not-found/not-found.component';
import {SignInComponent} from './account/sign-in/sign-in.component';
import {AuthUserActiveService} from './shared/route/auth-user-active.service';
import {PackageComponent} from './package/package.component';
import {UserComponent} from './user/user.component';
import {ClusterComponent} from './cluster/cluster.component';
import {ClusterDetailComponent} from './cluster/cluster-detail/cluster-detail.component';
import {OverviewComponent} from './overview/overview.component';
import {NodeComponent} from './node/node.component';
import {LogComponent} from './log/log.component';
import {ClusterRoutingResolverService} from './cluster/cluster-routing-resolver.service';
import {HostComponent} from './host/host.component';
import {DeployComponent} from './deploy/deploy.component';
import {SettingComponent} from './setting/setting.component';
import {SystemSettingComponent} from './setting/system-setting/system-setting.component';
import {CredentialComponent} from './credential/credential.component';
import {RegionComponent} from './region/region.component';
import {ZoneComponent} from './zone/zone.component';
import {PlanComponent} from './plan/plan.component';
import {F5BigIpComponent} from './f5-big-ip/f5-big-ip.component';
import {DeployPlanComponent} from './deploy-plan/deploy-plan.component';
import {ClusterHealthComponent} from './cluster-health/cluster-health.component';
import {ApplicationComponent} from './application/application.component';
import {ClusterBackupComponent} from './cluster-backup/cluster-backup.component';
import {BackupStorageSettingComponent} from './setting/backup-storage-setting/backup-storage-setting.component';
import {NfsComponent} from './nfs/nfs.component';
import {StorageComponent} from './storage/storage.component';
import {DashboardComponent} from './dashboard/dashboard.component';
import {SystemLogComponent} from './system-log/system-log.component';
import {DnsComponent} from './dns/dns.component';

const routes: Routes = [
  {path: '', redirectTo: 'kubeOperator', pathMatch: 'full'},
  {path: 'sign-in', component: SignInComponent},
  {
    path: 'kubeOperator',
    component: ShellComponent,
    canActivate: [AuthUserActiveService],
    canActivateChild: [AuthUserActiveService],
    children: [
      {path: '', redirectTo: 'cluster', pathMatch: 'full'},
      {path: 'dashboard', component: DashboardComponent},
      {path: 'cluster', component: ClusterComponent},
      {path: 'package', component: PackageComponent},
      {path: 'user', component: UserComponent},
      {path: 'host', component: HostComponent},
      {
        path: 'storage',
        component: StorageComponent,
        children: [
          {path: '', redirectTo: 'nfs', pathMatch: 'full'},
          {path: 'nfs', component: NfsComponent},
        ]
      },
      {
        path: 'plan',
        component: DeployPlanComponent,
        children: [
          {path: '', redirectTo: 'region', pathMatch: 'full'},
          {path: 'region', component: RegionComponent},
          {path: 'zone', component: ZoneComponent},
          {path: 'plan', component: PlanComponent}
        ]
      },
      {
        path: 'setting',
        component: SettingComponent,
        children: [
          {path: '', redirectTo: 'system', pathMatch: 'full'},
          {path: 'system', component: SystemSettingComponent},
          {path: 'credential', component: CredentialComponent},
          {path: 'backup-storage', component: BackupStorageSettingComponent},
          {path: 'dns', component: DnsComponent}
        ]
      },
      {
        path: 'cluster/:name',
        component: ClusterDetailComponent,
        resolve: {cluster: ClusterRoutingResolverService},
        children: [
          {path: '', redirectTo: 'overview', pathMatch: 'full'},
          {path: 'overview', component: OverviewComponent},
          {path: 'node', component: NodeComponent},
          {path: 'deploy', component: DeployComponent},
          {path: 'log', component: LogComponent},
          {path: 'apps', component: ApplicationComponent},
          {path: 'health', component: ClusterHealthComponent},
          {path: 'backup', component: ClusterBackupComponent},
          {path: 'big-ip', component: F5BigIpComponent}
        ]
      },
      {
        path: 'system/log', component: SystemLogComponent
      }
    ]
  },
  {path: '**', component: NotFoundComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
