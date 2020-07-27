import {NgModule} from '@angular/core';
import {ClusterComponent} from './cluster.component';
import {CoreModule} from '../../core/core.module';
import {ClusterListComponent} from './cluster-list/cluster-list.component';
import {ClusterCreateComponent} from './cluster-create/cluster-create.component';
import {ClusterDeleteComponent} from './cluster-delete/cluster-delete.component';
import {ClusterDetailComponent} from './cluster-detail/cluster-detail.component';
import {OverviewComponent} from './cluster-detail/overview/overview.component';
import {RouterModule} from '@angular/router';
import {ClusterConditionComponent} from './cluster-condition/cluster-condition.component';
import {NodeComponent} from './cluster-detail/node/node.component';
import {NamespaceComponent} from './cluster-detail/namespace/namespace.component';
import {NamespaceListComponent} from './cluster-detail/namespace/namespace-list/namespace-list.component';
import {SharedModule} from '../../shared/shared.module';
import {StorageComponent} from './cluster-detail/storage/storage.component';
import {PersistentVolumeComponent} from './cluster-detail/storage/persistent-volume/persistent-volume.component';
import {PersistentVolumeClaimComponent} from './cluster-detail/storage/persistent-volume-claim/persistent-volume-claim.component';
import {StorageClassComponent} from './cluster-detail/storage/storage-class/storage-class.component';
import {PersistentVolumeListComponent} from './cluster-detail/storage/persistent-volume/persistent-volume-list/persistent-volume-list.component';
import {NodeListComponent} from './cluster-detail/node/node-list/node-list.component';
import {NodeDetailComponent} from './cluster-detail/node/node-detail/node-detail.component';
import {PersistentVolumeClaimListComponent} from './cluster-detail/storage/persistent-volume-claim/persistent-volume-claim-list/persistent-volume-claim-list.component';
import {StorageClassListComponent} from './cluster-detail/storage/storage-class/storage-class-list/storage-class-list.component';
import {LoggingComponent} from './cluster-detail/logging/logging.component';
import {LoggingQueryComponent} from './cluster-detail/logging/logging-query/logging-query.component';

import {NgCircleProgressModule} from 'ng-circle-progress';
import {MonitorComponent} from './cluster-detail/monitor/monitor.component';
import {MonitorDashboardComponent} from './cluster-detail/monitor/monitor-dashboard/monitor-dashboard.component';
import {ClusterImportComponent} from './cluster-import/cluster-import.component';
import {StorageClassCreateComponent} from './cluster-detail/storage/storage-class/storage-class-create/storage-class-create.component';
import {PersistentVolumeCreateComponent} from './cluster-detail/storage/persistent-volume/persistent-volume-create/persistent-volume-create.component';
import {PersistentVolumeCreateNfsComponent} from './cluster-detail/storage/persistent-volume/persistent-volume-create/persistent-volume-create-nfs/persistent-volume-create-nfs.component';
import {PersistentVolumeCreateHostPathComponent} from './cluster-detail/storage/persistent-volume/persistent-volume-create/persistent-volume-create-host-path/persistent-volume-create-host-path.component';
import {StorageProvisionerComponent} from './cluster-detail/storage/storage-provisioner/storage-provisioner.component';
import {StorageProvisionerListComponent} from './cluster-detail/storage/storage-provisioner/storage-provisioner-list/storage-provisioner-list.component';
import {StorageProvisionerCreateComponent} from './cluster-detail/storage/storage-provisioner/storage-provisioner-create/storage-provisioner-create.component';
import {StorageProvisionerCreateNfsComponent} from './cluster-detail/storage/storage-provisioner/storage-provisioner-create/storage-provisioner-create-nfs/storage-provisioner-create-nfs.component';
import {StorageProvisionerDeleteComponent} from './cluster-detail/storage/storage-provisioner/storage-provisioner-delete/storage-provisioner-delete.component';
import {ToolsComponent} from './cluster-detail/tools/tools.component';
import {ToolsListComponent} from './cluster-detail/tools/tools-list/tools-list.component';
import {RepositoryComponent} from './cluster-detail/repository/repository.component';
import {RegistryComponent} from './cluster-detail/repository/registry/registry.component';
import {ChartmuseumComponent} from './cluster-detail/repository/chartmuseum/chartmuseum.component';
import {ChartListComponent} from './cluster-detail/repository/chartmuseum/chart-list/chart-list.component';
import {RegistryListComponent} from './cluster-detail/repository/registry/registry-list/registry-list.component';
import {ToolsEnableComponent} from './cluster-detail/tools/tools-enable/tools-enable.component';
import {NotReadyComponent} from './cluster-detail/not-ready/not-ready.component';
import { DashboardComponent } from './cluster-detail/dashboard/dashboard.component';
import { ToolsFailedComponent } from './cluster-detail/tools/tools-failed/tools-failed.component';
import { DashboardDashboardComponent } from './cluster-detail/dashboard/dashboard-dashboard/dashboard-dashboard.component';
import { NodeCreateComponent } from './cluster-detail/node/node-create/node-create.component';
import { NodeDeleteComponent } from './cluster-detail/node/node-delete/node-delete.component';
import { NodeStatusComponent } from './cluster-detail/node/node-status/node-status.component';
import { WebkubectlComponent } from './cluster-detail/overview/webkubectl/webkubectl.component';
import { CatalogComponent } from './cluster-detail/catalog/catalog.component';
import { CatalogDashboardComponent } from './cluster-detail/catalog/catalog-dashboard/catalog-dashboard.component';


@NgModule({
    declarations: [ClusterComponent, ClusterListComponent, ClusterCreateComponent, ClusterDeleteComponent, ClusterDetailComponent,
        OverviewComponent, ClusterConditionComponent, NodeComponent, NamespaceComponent, NamespaceListComponent,
        StorageComponent, PersistentVolumeComponent, PersistentVolumeClaimComponent, StorageClassComponent, PersistentVolumeListComponent,
        NodeListComponent, NodeDetailComponent, PersistentVolumeClaimListComponent,
        StorageClassListComponent, LoggingComponent, LoggingQueryComponent,
        MonitorComponent, MonitorDashboardComponent,
        ClusterImportComponent,
        StorageClassCreateComponent,
        PersistentVolumeCreateComponent,
        PersistentVolumeCreateNfsComponent,
        PersistentVolumeCreateHostPathComponent,
        StorageProvisionerComponent,
        StorageProvisionerListComponent,
        StorageProvisionerCreateComponent,
        StorageProvisionerCreateNfsComponent,
        StorageProvisionerDeleteComponent,
        ToolsComponent,
        ToolsListComponent,
        RepositoryComponent,
        RegistryComponent,
        ChartmuseumComponent,
        ChartListComponent,
        RegistryListComponent,
        ToolsEnableComponent,
        NotReadyComponent,
        DashboardComponent,
        ToolsFailedComponent,
        DashboardDashboardComponent,
        NodeCreateComponent,
        NodeDeleteComponent,
        NodeStatusComponent,
        WebkubectlComponent,
        CatalogComponent,
        CatalogDashboardComponent],
    imports: [
        CoreModule,
        RouterModule,
        SharedModule,
        NgCircleProgressModule,
    ]
})
export class ClusterModule {
}
