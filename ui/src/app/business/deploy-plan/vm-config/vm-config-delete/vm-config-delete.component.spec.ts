import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VmConfigDeleteComponent } from './vm-config-delete.component';

describe('VmConfigDeleteComponent', () => {
  let component: VmConfigDeleteComponent;
  let fixture: ComponentFixture<VmConfigDeleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ VmConfigDeleteComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(VmConfigDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
