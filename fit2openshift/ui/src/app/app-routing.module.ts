import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {ShellComponent} from './base/shell/shell.component';
import {NotFoundComponent} from './shared/not-found/not-found.component';
import {SignInComponent} from './account/sign-in/sign-in.component';
import {AuthUserActiveService} from './shared/route/auth-user-active.service';
import {OfflineComponent} from './offline/offline.component';
import {UserComponent} from './user/user.component';
import {ClusterComponent} from './cluster/cluster.component';
import {ClusterDetailComponent} from './cluster/cluster-detail/cluster-detail.component';
import {OverviewComponent} from './overview/overview.component';
import {NodeComponent} from './node/node.component';
import {LogComponent} from './log/log.component';
import {ConfigComponent} from './config/config.component';
import {MonitorComponent} from './monitor/monitor.component';

const routes: Routes = [
  {path: '', redirectTo: 'fit2openshift', pathMatch: 'full'},
  {path: 'sign-in', component: SignInComponent},
  {
    path: 'fit2openshift',
    component: ShellComponent,
    // canActivate: [AuthUserActiveService],
    // canActivateChild: [AuthUserActiveService],
    children: [
      {path: '', redirectTo: 'offline', pathMatch: 'full'},
      {path: 'cluster', component: ClusterComponent},
      {path: 'offline', component: OfflineComponent},
      {path: 'user', component: UserComponent},
      {
        path: 'cluster/:id',
        component: ClusterDetailComponent,
        children: [
          {path: 'overview', component: OverviewComponent},
          {path: 'node', component: NodeComponent},
          {path: 'log', component: LogComponent},
          {path: 'config', component: ConfigComponent},
          {path: 'monitor', component: MonitorComponent}
        ]
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
