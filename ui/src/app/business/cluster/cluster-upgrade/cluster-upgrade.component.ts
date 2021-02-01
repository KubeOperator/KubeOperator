import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Cluster, ClusterUpgradeRequest} from '../cluster';
import {ManifestService} from '../../manifest/manifest.service';
import {Manifest, NameVersion} from '../../manifest/manifest';
import {NgForm} from '@angular/forms';
import {ClusterService} from '../cluster.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {TranslateService} from '@ngx-translate/core';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';

@Component({
    selector: 'app-cluster-upgrade',
    templateUrl: './cluster-upgrade.component.html',
    styleUrls: ['./cluster-upgrade.component.css']
})
export class ClusterUpgradeComponent implements OnInit {

    opened = false;
    currentCluster: Cluster;
    isSubmitGoing = false;
    @Output() upgrade = new EventEmitter();
    upgradeVersions: string[] = [];
    @ViewChild('clusterForm') clusterForm: NgForm;
    clusterUpgradeRequest: ClusterUpgradeRequest = new ClusterUpgradeRequest();
    chooseVersion: string;
    oldManifest: Manifest;
    newManifest: Manifest;
    manifests: Manifest[] = [];


    constructor(private manifestService: ManifestService,
                private clusterService: ClusterService,
                private translateService: TranslateService,
                private commonAlertService: CommonAlertService) {
    }

    ngOnInit(): void {
    }

    open(item: Cluster) {
        this.opened = true;
        this.currentCluster = item;
        const currentVersion = this.currentCluster.spec.version;
        const currentVersions = currentVersion.split('.');
        const version1 = currentVersions[0];
        const version2 = currentVersions[1];
        const versions = currentVersions[2].split('-ko');
        const version3 = Number(versions[0]);
        const koVersion = Number(versions[1]);
        this.manifestService.listActive().subscribe(res => {
            this.manifests = res;
            for (const manifest of res) {
                const manifestKoVersions = manifest.name.split('-ko');
                const manifestVersions = manifestKoVersions[0].split('.');
                const manifestKoVersion = Number(manifestKoVersions[1]);
                const manifestVersion1 = manifestVersions[0];
                const manifestVersion2 = manifestVersions[1];
                const manifestVersion3 = Number(manifestVersions[2]);
                if (version1 === manifestVersion1 && version2 === manifestVersion2) {
                    if (manifestVersion3 > version3) {
                        this.upgradeVersions.push(manifest.name);
                    }
                    if (manifestVersion3 === version3 && koVersion < manifestKoVersion) {
                        this.upgradeVersions.push(manifest.name);
                    }
                }
            }
        });
    }

    onCancel() {
        this.opened = false;
        this.upgradeVersions = [];
        this.oldManifest = null;
        this.newManifest = null;
        this.clusterUpgradeRequest = new ClusterUpgradeRequest();
        this.clusterForm.resetForm(this.upgradeVersions);
        this.clusterForm.resetForm(this.currentCluster);
    }


    onSelectChooseVersion() {
        for (const m of this.manifests) {
            if (m.name.indexOf(this.currentCluster.spec.version) !== -1) {
                this.oldManifest = m;
            }
            if (m.name === this.chooseVersion) {
                this.newManifest = m;
            }
        }
    }

    onSubmit() {
        this.clusterService.upgrade(this.currentCluster.name, this.chooseVersion).subscribe(res => {
            this.onCancel();
            this.upgrade.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPGRADE_START_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.onCancel();
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    getVersion(component: string, ns: NameVersion[]): string {
        for (const n of ns) {
            if (n.name === component) {
                return n.version;
            }
        }
    }
}
