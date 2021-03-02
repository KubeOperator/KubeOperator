import {BaseModelService} from './BaseModelService';
import {Directive, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModel} from './BaseModel';
import {MultiClusterRepository} from '../../business/multi-cluster/multi-cluster-repository';

@Directive()
export abstract class BaseModelDirective<T extends BaseModel> implements OnInit {

    items: T[] = [];
    page = 1;
    size = 10;
    total = 0;
    loading = true;
    selected: T[] = [];
    @Output() createEvent = new EventEmitter();
    @Output() deleteEvent = new EventEmitter<T[]>();
    @Output() updateEvent = new EventEmitter<T>();
    @Output() errorDetailEvent = new EventEmitter<string>();

    protected constructor(protected service: BaseModelService<T>) {
    }

    ngOnInit(): void {
        this.refresh();
    }

    onCreate() {
        this.createEvent.emit();
    }

    onDelete() {
        this.deleteEvent.emit(this.selected);
    }

    onErrorDetail(item: MultiClusterRepository) {
        this.errorDetailEvent.emit(item.message);
    }

    reset() {
        this.selected = [];
    }

    onUpdate(item: T) {
        this.updateEvent.emit(item);
    }

    refresh() {
        this.loading = true;
        this.service.page(this.page, this.size).subscribe(data => {
            this.items = data.items;
            this.total = data.total;
            this.loading = false;
        });
    }
}
