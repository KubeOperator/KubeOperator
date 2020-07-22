import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {LoginComponent} from './login/login.component';
import {LayoutComponent} from './layout/layout.component';
import {ClusterComponent} from './business/cluster/cluster.component';
import {ClusterDetailComponent} from './business/cluster/cluster-detail/cluster-detail.component';
import {ClusterRoutingResolverService} from './business/cluster/cluster-routing-resolver.service';
import {OverviewComponent} from './business/cluster/cluster-detail/overview/overview.component';
import {SettingComponent} from './business/setting/setting.component';
import {CredentialComponent} from './business/setting/credential/credential.component';
import {HostComponent} from './business/host/host.component';
import {NodeComponent} from './business/cluster/cluster-detail/node/node.component';
import {NamespaceComponent} from './business/cluster/cluster-detail/namespace/namespace.component';
import {StorageComponent} from './business/cluster/cluster-detail/storage/storage.component';
import {PersistentVolumeComponent} from './business/cluster/cluster-detail/storage/persistent-volume/persistent-volume.component';
import {PersistentVolumeClaimComponent} from './business/cluster/cluster-detail/storage/persistent-volume-claim/persistent-volume-claim.component';
import {UserComponent} from './business/user/user.component';
import {AuthUserService} from './shared/auth/auth-user.service';
import {LoggingComponent} from './business/cluster/cluster-detail/logging/logging.component';
import {MonitorComponent} from './business/cluster/cluster-detail/monitor/monitor.component';
import {StorageClassComponent} from './business/cluster/cluster-detail/storage/storage-class/storage-class.component';
import {RegionComponent} from './business/deploy-plan/region/region.component';
import {DeployPlanComponent} from './business/deploy-plan/deploy-plan.component';
import {ZoneComponent} from './business/deploy-plan/zone/zone.component';
import {PlanComponent} from './business/deploy-plan/plan/plan.component';
import {StorageProvisionerComponent} from './business/cluster/cluster-detail/storage/storage-provisioner/storage-provisioner.component';
import {RepositoryComponent} from './business/cluster/cluster-detail/repository/repository.component';
import {ChartmuseumComponent} from './business/cluster/cluster-detail/repository/chartmuseum/chartmuseum.component';
import {RegistryComponent} from './business/cluster/cluster-detail/repository/registry/registry.component';
import {ToolsComponent} from './business/cluster/cluster-detail/tools/tools.component';
import {DashboardComponent} from './business/cluster/cluster-detail/dashboard/dashboard.component';
import {SystemComponent} from './business/setting/system/system.component';
import {ProjectComponent} from './business/project/project.component';
import {ProjectDetailComponent} from './business/project/project-detail/project-detail.component';
import {ProjectRoutingResolverService} from './business/project/project-routing-resolver.service';
import {ProjectResourceComponent} from './business/project/project-resource/project-resource.component';

const routes: Routes = [
    {path: 'login', component: LoginComponent},
    {
        path: '',
        component: LayoutComponent,
        canActivate: [AuthUserService],
        canActivateChild: [AuthUserService],
        children: [
            {path: '', redirectTo: 'clusters', pathMatch: 'full'},
            {
                path: 'projects',
                component: ProjectComponent,
            },
            {
                path: 'projects/:name',
                component: ProjectDetailComponent,
                resolve: {project: ProjectRoutingResolverService},
                children: [
                    {path: '', redirectTo: 'clusters', pathMatch: 'full'},
                    {path: 'clusters', component: ClusterComponent},
                    {path: 'resources', component: ProjectResourceComponent}
                ]
            },
            {
                path: 'clusters',
                component: ClusterComponent,
            },
            {
                path: 'clusters/:name',
                component: ClusterDetailComponent,
                resolve: {cluster: ClusterRoutingResolverService},
                children: [
                    {path: '', redirectTo: 'overview', pathMatch: 'full'},
                    {path: 'overview', component: OverviewComponent},
                    {path: 'nodes', component: NodeComponent},
                    {path: 'namespaces', component: NamespaceComponent},
                    {
                        path: 'storages',
                        component: StorageComponent,
                        children: [
                            {path: '', redirectTo: 'pv', pathMatch: 'full'},
                            {path: 'pv', component: PersistentVolumeComponent},
                            {path: 'pvc', component: PersistentVolumeClaimComponent},
                            {path: 'sc', component: StorageClassComponent},
                            {path: 'provisioner', component: StorageProvisionerComponent},
                        ],
                    },
                    {path: 'logging', component: LoggingComponent},
                    {path: 'monitor', component: MonitorComponent},
                    {
                        path: 'repository',
                        component: RepositoryComponent,
                        children: [
                            {path: '', redirectTo: 'chartmuseum', pathMatch: 'full'},
                            {path: 'chartmuseum', component: ChartmuseumComponent},
                            {path: 'registry', component: RegistryComponent}
                        ]
                    },
                    {path: 'tool', component: ToolsComponent},
                    {path: 'dashboard', component: DashboardComponent},
                ],
            },
            {
                path: 'hosts',
                component: HostComponent,
            },
            {
                path: 'setting',
                component: SettingComponent,
                children: [
                    {path: '', redirectTo: 'system', pathMatch: 'full'},
                    {path: 'system', component: SystemComponent},
                    {path: 'credential', component: CredentialComponent},
                ]
            },
            {
                path: 'deploy',
                component: DeployPlanComponent,
                children: [
                    {path: '', redirectTo: 'region', pathMatch: 'full'},
                    {path: 'region', component: RegionComponent},
                    {path: 'zone', component: ZoneComponent},
                    {path: 'plan', component: PlanComponent}
                ]
            },
            {
                path: 'users',
                component: UserComponent,
            }
        ]
    }
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutingModule {
}
