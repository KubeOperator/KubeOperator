import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PersistentVolumeCreateComponent } from './persistent-volume-create.component';

describe('PersistentVolumeCreateComponent', () => {
  let component: PersistentVolumeCreateComponent;
  let fixture: ComponentFixture<PersistentVolumeCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PersistentVolumeCreateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PersistentVolumeCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
