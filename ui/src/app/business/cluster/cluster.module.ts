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
import {EfComponent} from './cluster-detail/logging/ef/ef.component';
import {LokiComponent} from './cluster-detail/logging/loki/loki.component';
import {IstioComponent} from './cluster-detail/istio/istio.component';

import {NgCircleProgressModule} from 'ng-circle-progress';
import {MonitorComponent} from './cluster-detail/monitor/monitor.component';
import {MonitorDashboardComponent} from './cluster-detail/monitor/monitor-dashboard/monitor-dashboard.component';
import {ClusterImportComponent} from './cluster-import/cluster-import.component';
import {StorageClassCreateComponent} from './cluster-detail/storage/storage-class/storage-class-create/storage-class-create.component';
import {PersistentVolumeCreateComponent} from './cluster-detail/storage/persistent-volume/persistent-volume-create/persistent-volume-create.component';
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
import { ToolsFailedComponent } from './cluster-detail/tools/tools-failed/tools-failed.component';
import { NodeCreateComponent } from './cluster-detail/node/node-create/node-create.component';
import { NodeDeleteComponent } from './cluster-detail/node/node-delete/node-delete.component';
import { NodeStatusComponent } from './cluster-detail/node/node-status/node-status.component';
import { WebkubectlComponent } from './cluster-detail/overview/webkubectl/webkubectl.component';
import { StorageProvisionerCreateExternalCephComponent } from './cluster-detail/storage/storage-provisioner/storage-provisioner-create/storage-provisioner-create-external-ceph/storage-provisioner-create-external-ceph.component';
import { StorageProvisionerCreateRookCephComponent } from './cluster-detail/storage/storage-provisioner/storage-provisioner-create/storage-provisioner-create-rook-ceph/storage-provisioner-create-rook-ceph.component';
import { StorageProvisionerCreateVsphereComponent } from './cluster-detail/storage/storage-provisioner/storage-provisioner-create/storage-provisioner-create-vsphere/storage-provisioner-create-vsphere.component';
import { BackupComponent } from './cluster-detail/backup/backup.component';
import {BackupModule} from './cluster-detail/backup/backup.module';
import { ToolsDisableComponent } from './cluster-detail/tools/tools-disable/tools-disable.component';
import { PersistentVolumeCreateLocalStorageComponent } from './cluster-detail/storage/persistent-volume/persistent-volume-create/persistent-volume-create-local-storage/persistent-volume-create-local-storage.component';
import { SecurityComponent } from './cluster-detail/security/security.component';
import { SecurityTaskListComponent } from './cluster-detail/security/security-task-list/security-task-list.component';
import { SecurityTaskCreateComponent } from './cluster-detail/security/security-task-create/security-task-create.component';
import { SecurityTaskDetailComponent } from './cluster-detail/security/security-task-detail/security-task-detail.component';
import { SecurityTaskDeleteComponent } from './cluster-detail/security/security-task-delete/security-task-delete.component';
import { EventComponent } from './cluster-detail/event/event.component';
import { ClusterUpgradeComponent } from './cluster-upgrade/cluster-upgrade.component';
import { ClusterLoggerComponent } from './cluster-logger/cluster-logger.component';
import { ClusterGradeComponent } from './cluster-detail/cluster-grade/cluster-grade.component';
import {NgxEchartsModule} from 'ngx-echarts';
import { NamespaceDeleteComponent } from './cluster-detail/namespace/namespace-delete/namespace-delete.component';
import { NamespaceCreateComponent } from './cluster-detail/namespace/namespace-create/namespace-create.component';
import { StorageProvisionerCreateOceanStorComponent } from './cluster-detail/storage/storage-provisioner/storage-provisioner-create/storage-provisioner-create-ocean-stor/storage-provisioner-create-ocean-stor.component';
import { F5Component } from './cluster-detail/f5/f5.component';
import { ClusterHealthCheckComponent } from './cluster-health-check/cluster-health-check.component';


@NgModule({
    declarations: [ClusterComponent, ClusterListComponent, ClusterCreateComponent, ClusterDeleteComponent, ClusterDetailComponent,
        OverviewComponent, ClusterConditionComponent, NodeComponent, NamespaceComponent, NamespaceListComponent,
        StorageComponent, PersistentVolumeComponent, PersistentVolumeClaimComponent, StorageClassComponent, PersistentVolumeListComponent,
        NodeListComponent, NodeDetailComponent, PersistentVolumeClaimListComponent,
        StorageClassListComponent, LoggingComponent, EfComponent, LokiComponent, IstioComponent,
        MonitorComponent, MonitorDashboardComponent,
        ClusterImportComponent,
        StorageClassCreateComponent,
        PersistentVolumeCreateComponent,
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
        ToolsFailedComponent,
        NodeCreateComponent,
        NodeDeleteComponent,
        NodeStatusComponent,
        WebkubectlComponent,
        StorageProvisionerCreateExternalCephComponent,
        StorageProvisionerCreateRookCephComponent,
        StorageProvisionerCreateVsphereComponent,
        BackupComponent,
        ToolsDisableComponent,
        PersistentVolumeCreateLocalStorageComponent,
        SecurityComponent,
        SecurityTaskListComponent,
        SecurityTaskCreateComponent,
        SecurityTaskDetailComponent,
        SecurityTaskDeleteComponent,
        EventComponent,
        ClusterUpgradeComponent,
        ClusterLoggerComponent,
        ClusterGradeComponent,
        NamespaceDeleteComponent,
        NamespaceCreateComponent,
        StorageProvisionerCreateOceanStorComponent,
        F5Component,
        ClusterHealthCheckComponent],
    imports: [
        CoreModule,
        RouterModule,
        SharedModule,
        NgCircleProgressModule,
        BackupModule,
        NgxEchartsModule.forRoot({
            echarts: () => import('echarts'),
        }),
    ]
})
export class ClusterModule {
}
