import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PersistentVolumeCreateLocalStorageComponent } from './persistent-volume-create-local-storage.component';

describe('PersistentVolumeCreateLocalStorageComponent', () => {
  let component: PersistentVolumeCreateLocalStorageComponent;
  let fixture: ComponentFixture<PersistentVolumeCreateLocalStorageComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PersistentVolumeCreateLocalStorageComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PersistentVolumeCreateLocalStorageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
