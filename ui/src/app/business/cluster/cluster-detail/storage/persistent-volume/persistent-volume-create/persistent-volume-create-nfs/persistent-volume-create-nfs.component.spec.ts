import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PersistentVolumeCreateNfsComponent } from './persistent-volume-create-nfs.component';

describe('PersistentVolumeCreateNfsComponent', () => {
  let component: PersistentVolumeCreateNfsComponent;
  let fixture: ComponentFixture<PersistentVolumeCreateNfsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PersistentVolumeCreateNfsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PersistentVolumeCreateNfsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
