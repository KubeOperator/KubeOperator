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
import {ServiceComponent} from './cluster-detail/service/service.component';
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
import {ConfigComponent} from './cluster-detail/config/config.component';
import {PersistentVolumeClaimListComponent} from './cluster-detail/storage/persistent-volume-claim/persistent-volume-claim-list/persistent-volume-claim-list.component';
import {StorageClassListComponent} from './cluster-detail/storage/storage-class/storage-class-list/storage-class-list.component';
import {WorkloadComponent} from './cluster-detail/workload/workload.component';
import {DeploymentComponent} from './cluster-detail/workload/deployment/deployment.component';
import {StatefulSetComponent} from './cluster-detail/workload/stateful-set/stateful-set.component';
import {DaemonSetComponent} from './cluster-detail/workload/daemon-set/daemon-set.component';
import {JobComponent} from './cluster-detail/workload/job/job.component';
import {CornJobComponent} from './cluster-detail/workload/corn-job/corn-job.component';
import {DeploymentListComponent} from './cluster-detail/workload/deployment/deployment-list/deployment-list.component';
import {StatefulSetListComponent} from './cluster-detail/workload/stateful-set/stateful-set-list/stateful-set-list.component';
import {JobListComponent} from './cluster-detail/workload/job/job-list/job-list.component';
import {CornJobListComponent} from './cluster-detail/workload/corn-job/corn-job-list/corn-job-list.component';
import {DaemonSetListComponent} from './cluster-detail/workload/daemon-set/daemon-set-list/daemon-set-list.component';
import {ServiceListComponent} from './cluster-detail/service/service-list/service-list.component';
import {ConfigMapComponent} from './cluster-detail/config/config-map/config-map.component';
import {SecretComponent} from './cluster-detail/config/secret/secret.component';
import {ConfigMapListComponent} from './cluster-detail/config/config-map/config-map-list/config-map-list.component';
import {SecretListComponent} from './cluster-detail/config/secret/secret-list/secret-list.component';
import {LoggingComponent} from './cluster-detail/logging/logging.component';
import {LoggingQueryComponent} from './cluster-detail/logging/logging-query/logging-query.component';

import {NgCircleProgressModule} from 'ng-circle-progress';
import {MonitorComponent} from './cluster-detail/monitor/monitor.component';
import {MonitorDashboardComponent} from './cluster-detail/monitor/monitor-dashboard/monitor-dashboard.component';
import {MonitorEnableComponent} from './cluster-detail/monitor/monitor-enable/monitor-enable.component';
import {MonitorStatusComponent} from './cluster-detail/monitor/monitor-status/monitor-status.component';
import { ClusterImportComponent } from './cluster-import/cluster-import.component';
import { StorageClassCreateComponent } from './cluster-detail/storage/storage-class/storage-class-create/storage-class-create.component';
import { PersistentVolumeCreateComponent } from './cluster-detail/storage/persistent-volume/persistent-volume-create/persistent-volume-create.component';
import { PersistentVolumeCreateNfsComponent } from './cluster-detail/storage/persistent-volume/persistent-volume-create/persistent-volume-create-nfs/persistent-volume-create-nfs.component';
import { PersistentVolumeCreateHostPathComponent } from './cluster-detail/storage/persistent-volume/persistent-volume-create/persistent-volume-create-host-path/persistent-volume-create-host-path.component';
import { StorageProvisionerComponent } from './cluster-detail/storage/storage-provisioner/storage-provisioner.component';
import { StorageProvisionerListComponent } from './cluster-detail/storage/storage-provisioner/storage-provisioner-list/storage-provisioner-list.component';
import { StorageProvisionerCreateComponent } from './cluster-detail/storage/storage-provisioner/storage-provisioner-create/storage-provisioner-create.component';
import { StorageProvisionerCreateNfsComponent } from './cluster-detail/storage/storage-provisioner/storage-provisioner-create/storage-provisioner-create-nfs/storage-provisioner-create-nfs.component';
import { StorageProvisionerDeleteComponent } from './cluster-detail/storage/storage-provisioner/storage-provisioner-delete/storage-provisioner-delete.component';
import { ToolsComponent } from './cluster-detail/tools/tools.component';
import { ToolsListComponent } from './cluster-detail/tools/tools-list/tools-list.component';
import { RepositoryComponent } from './cluster-detail/repository/repository.component';
import { RegistryComponent } from './cluster-detail/repository/registry/registry.component';
import { ChartmuseumComponent } from './cluster-detail/repository/chartmuseum/chartmuseum.component';
import { ChartListComponent } from './cluster-detail/repository/chartmuseum/chart-list/chart-list.component';
import { RegistryListComponent } from './cluster-detail/repository/registry/registry-list/registry-list.component';
import { ToolsEnableComponent } from './cluster-detail/tools/tools-enable/tools-enable.component';


@NgModule({
    declarations: [ClusterComponent, ClusterListComponent, ClusterCreateComponent, ClusterDeleteComponent, ClusterDetailComponent,
        OverviewComponent, ClusterConditionComponent, NodeComponent, NamespaceComponent, NamespaceListComponent,
        StorageComponent, PersistentVolumeComponent, PersistentVolumeClaimComponent, StorageClassComponent, PersistentVolumeListComponent,
        NodeListComponent, NodeDetailComponent, ConfigComponent, PersistentVolumeClaimListComponent,
        StorageClassListComponent, WorkloadComponent, DeploymentComponent, StatefulSetComponent, DaemonSetComponent, JobComponent,
        CornJobComponent, DeploymentListComponent, StatefulSetListComponent, JobListComponent, CornJobListComponent,
        DaemonSetListComponent, ServiceComponent,  ServiceListComponent,  ConfigMapComponent,
        SecretComponent, ConfigMapListComponent, SecretListComponent, LoggingComponent, LoggingQueryComponent,
        MonitorComponent, MonitorDashboardComponent, MonitorEnableComponent,
        MonitorStatusComponent,
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
        ToolsEnableComponent],
    imports: [
        CoreModule,
        RouterModule,
        SharedModule,
        NgCircleProgressModule,
    ]
})
export class ClusterModule {
}
