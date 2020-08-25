import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Category, Manifest} from "../manifest";
import {ManifestService} from "../manifest.service";

@Component({
    selector: 'app-manifest-list',
    templateUrl: './manifest-list.component.html',
    styleUrls: ['./manifest-list.component.css']
})
export class ManifestListComponent implements OnInit {
    constructor(private manifestService: ManifestService) {
    }

    items: Manifest[] = [];
    loading = false;
    @Output() detailEvent = new EventEmitter<Manifest>();

    ngOnInit(): void {
        this.refresh();
    }

    refresh() {
        this.manifestService.list().subscribe(data => {
            this.items = data;
            console.log(this.items);
        });
    }

    onDetail(item: Manifest) {
        this.detailEvent.emit(item);
    }

    getCoreCateGory(item: Manifest): Category {
        for (const category of item.category) {
            if (category.name === 'core') {
                return category;
            }
        }
    }

}
