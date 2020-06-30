import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PersistentVolumeCreateHostPathComponent } from './persistent-volume-create-host-path.component';

describe('PersistentVolumeCreateHostPathComponent', () => {
  let component: PersistentVolumeCreateHostPathComponent;
  let fixture: ComponentFixture<PersistentVolumeCreateHostPathComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PersistentVolumeCreateHostPathComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PersistentVolumeCreateHostPathComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
